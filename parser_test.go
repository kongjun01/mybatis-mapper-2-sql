package parser

import (
	"testing"
)

func testParser(t *testing.T, xmlData, expect string) {
	actual, err := ParseXML(xmlData)
	if err != nil {
		t.Errorf("parse error: %v", err)
		return
	}
	if actual != expect {
		t.Errorf("\nexpect: [%s]\nactual: [%s]", expect, actual)
	}
}

func TestParserIf(t *testing.T) {
	testParser(t,
		`
<mapper namespace="Test">
	   <select id="testIf">
        SELECT
        name,
        category,
        price
        FROM
        fruits
        WHERE
        1=1
        <if test="category != null and category !=''">
            AND category = #{category}
        </if>
        <if test="price != null and price !=''">
            AND price = ${price}
            <if test="price >= 400">
                AND name = 'Fuji'
            </if>
        </if>
    </select>
</mapper>`,
		"SELECT `name`,`category`,`price` FROM `fruits` WHERE 1=1 AND `category`=? AND `price`=? AND `name`=\"Fuji\";",
	)
}

func TestParserParams(t *testing.T) {
	testParser(t,
		`
<mapper namespace="Test">
    <select id="testParameters">
        SELECT
        name,
        category,
        price
        FROM
        fruits
        WHERE
        category = #{category}
        AND price > ${price}
    </select>
</mapper>`,
		"SELECT `name`,`category`,`price` FROM `fruits` WHERE `category`=? AND `price`>?;",
	)
}

func TestParserInclude(t *testing.T) {
	testParser(t,
		`
<mapper namespace="Test">
	<sql id="sometable">
  		${prefix}Table
	</sql>
	<sql id="someinclude">
  		from
    	<include refid="${include_target}"/>
	</sql>
	<select id="select" resultType="map">
		select
		field1, field2, field3
  		<include refid="someinclude">
    		<property name="prefix" value="Some"/>
    		<property name="include_target" value="sometable"/>
  		</include>
	</select>
</mapper>`,
		"SELECT `field1`,`field2`,`field3` FROM `SomeTable`;",
	)
}

func TestParserTrim(t *testing.T) {
	testParser(t,
		`
<mapper namespace="Test">
<select id="testTrim">
        SELECT
        name,
        category,
        price
        FROM
        fruits
        <trim prefix="WHERE" prefixOverrides="AND |OR ">
            OR category = 'apple'
            OR price = 200
        </trim>
    </select>
</mapper>`,
		"SELECT `name`,`category`,`price` FROM `fruits` WHERE `category`=\"apple\" OR `price`=200;",
	)
	testParser(t,
		`
<mapper namespace="Test">
<select id="testTrim">
       SELECT
       name,
       category,
       price
       FROM
       fruits
       <trim prefix="WHERE" prefixOverrides="AND |OR ">
           AND category = 'apple'
           OR price = 200
       </trim>
   </select>
</mapper>`,
		"SELECT `name`,`category`,`price` FROM `fruits` WHERE `category`=\"apple\" OR `price`=200;",
	)
	testParser(t,
		`
<mapper namespace="Test">
<select id="testTrim">
       SELECT
       name,
       category,
       price
       FROM
       fruits
       <where>
           AND category = 'apple'
           OR price = 200
       </where>
   </select>
</mapper>`,
		"SELECT `name`,`category`,`price` FROM `fruits` WHERE `category`=\"apple\" OR `price`=200;",
	)
	testParser(t,
		`
<mapper namespace="Test">
<select id="testTrim">
       SELECT
       name,
       category,
       price
       FROM
       fruits
       <where>
           OR category = 'apple'
           OR price = 200
       </where>
   </select>
</mapper>`,
		"SELECT `name`,`category`,`price` FROM `fruits` WHERE `category`=\"apple\" OR `price`=200;",
	)
}

func TestParserWhereAndIf(t *testing.T) {
	testParser(t,
		`
<mapper namespace="Test">
    <select id="testWhereIf">
        SELECT
        name,
        category,
        price
        FROM
        fruits
        <where>
            AND category = 'apple'
            <if test="price != null and price !=''">
                AND price = ${price}
            </if>
        </where>
    </select>
</mapper>`,
		"SELECT `name`,`category`,`price` FROM `fruits` WHERE `category`=\"apple\" AND `price`=?;",
	)
}

func TestParserSet(t *testing.T) {
	testParser(t,
		`
<mapper namespace="Test">
    <update id="testSet">
        UPDATE
        fruits
        <set>
            <if test="category != null and category !=''">
                category = #{category},
            </if>
            <if test="price != null and price !=''">
                price = ${price},
            </if>
        </set>
        WHERE
        name = #{name}
    </update>
</mapper>`,
		"UPDATE `fruits` SET `category`=?, `price`=? WHERE `name`=?;",
	)
}

func TestParserChoose(t *testing.T) {
	testParser(t,
		`
<mapper namespace="Test">
    <select id="testChoose">
        SELECT
        name,
        category,
        price
        FROM
        fruits
        <where>
            <choose>
                <when test="name != null">
                    AND name = #{name}
                </when>
                <when test="category == 'banana'">
                    AND category = #{category}
                    <if test="price != null and price !=''">
                        AND price = ${price}
                    </if>
                </when>
                <otherwise>
                    AND category = 'apple'
                </otherwise>
            </choose>
        </where>
    </select>
</mapper>`,
		"SELECT `name`,`category`,`price` FROM `fruits` WHERE `name`=? AND `category`=? AND `price`=? AND `category`=\"apple\";",
	)
}

func TestParserForeach(t *testing.T) {
	testParser(t,
		`
<mapper namespace="Test">
    <select id="testForeach">
        SELECT
        name,
        category,
        price
        FROM
        fruits
        <where>
            category = 'apple' AND
            <foreach collection="apples" item="name" open="(" close=")" separator="OR">
				name = #{name}
            </foreach>
        </where>
    </select>
</mapper>`,
		"SELECT `name`,`category`,`price` FROM `fruits` WHERE `category`=\"apple\" AND (`name`=? OR `name`=?);",
	)
	testParser(t,
		`
<mapper namespace="Test">
    <select id="testForeach">
        SELECT
        name,
        category,
        price
        FROM
        fruits
        <where>
            category = 'apple' AND
            <foreach collection="apples" item="name" open="(" close=")" separator="OR">
                <if test="name == 'Jonathan' or name == 'Fuji'">
                    name = #{name}
                </if>
                <if test="name == 'Jonathan' or name == 'Fuji'">
                    name like #{name}
                </if>
            </foreach>
        </where>
    </select>
</mapper>`,
		"SELECT `name`,`category`,`price` FROM `fruits` WHERE `category`=\"apple\" AND (`name`=? OR `name` LIKE ?);",
	)
	testParser(t,
		`
<mapper namespace="Test">
	<insert id="testInsertMulti">
        INSERT INTO
        fruits
        (
        name,
        category,
        price
        )
        VALUES
        <foreach collection="fruits" item="fruit" separator=",">
            (
            #{fruit.name},
            #{fruit.category},
            ${fruit.price}
            )
        </foreach>
    </insert>
</mapper>`,
		"INSERT INTO `fruits` (`name`,`category`,`price`) VALUES (?,?,?),(?,?,?);",
	)
}

func TestParserBind(t *testing.T) {
	testParser(t,
		`
<mapper namespace="Test">
    <select id="testBind">
        <bind name="likeName" value="'%' + name + '%'"/>
        SELECT
        name,
        category,
        price
        FROM
        fruits
        WHERE
        name like #{likeName}
    </select>
</mapper>`,
		"SELECT `name`,`category`,`price` FROM `fruits` WHERE `name` LIKE ?;",
	)
}

func TestParserFullFile(t *testing.T) {
	testParser(t,
		`
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE mapper PUBLIC "-//mybatis.org//DTD Mapper 3.0//EN" "http://mybatis.org/dtd/mybatis-3-mapper.dtd">
<mapper namespace="Test">
    <sql id="sometable">
        fruits
    </sql>
    <sql id="somewhere">
        WHERE
        category = #{category}
    </sql>
    <sql id="someinclude">
        FROM
        <include refid="${include_target}"/>
        <include refid="somewhere"/>
    </sql>
    <select id="testParameters">
        SELECT
        name,
        category,
        price
        FROM
        fruits
        WHERE
        category = #{category}
        AND price > ${price}
    </select>
    <select id="testInclude">
        SELECT
        name,
        category,
        price
        <include refid="someinclude">
            <property name="prefix" value="Some"/>
            <property name="include_target" value="sometable"/>
        </include>
    </select>
    <select id="testIf">
        SELECT
        name,
        category,
        price
        FROM
        fruits
        WHERE
        1=1
        <if test="category != null and category !=''">
            AND category = #{category}
        </if>
        <if test="price != null and price !=''">
            AND price = ${price}
            <if test="price >= 400">
                AND name = 'Fuji'
            </if>
        </if>
    </select>
    <select id="testTrim">
        SELECT
        name,
        category,
        price
        FROM
        fruits
        <trim prefix="WHERE" prefixOverrides="AND|OR">
            OR category = 'apple'
            OR price = 200
        </trim>
    </select>
    <select id="testWhere">
        SELECT
        name,
        category,
        price
        FROM
        fruits
        <where>
            AND category = 'apple'
            <if test="price != null and price !=''">
                AND price = ${price}
            </if>
        </where>
    </select>
    <update id="testSet">
        UPDATE
        fruits
        <set>
            <if test="category != null and category !=''">
                category = #{category},
            </if>
            <if test="price != null and price !=''">
                price = ${price},
            </if>
        </set>
        WHERE
        name = #{name}
    </update>
    <select id="testChoose">
        SELECT
        name,
        category,
        price
        FROM
        fruits
        <where>
            <choose>
                <when test="name != null">
                    AND name = #{name}
                </when>
                <when test="category == 'banana'">
                    AND category = #{category}
                    <if test="price != null and price !=''">
                        AND price = ${price}
                    </if>
                </when>
                <otherwise>
                    AND category = 'apple'
                </otherwise>
            </choose>
        </where>
    </select>
    <select id="testForeach">
        SELECT
        name,
        category,
        price
        FROM
        fruits
        <where>
            category = 'apple' AND
            <foreach collection="apples" item="name" open="(" close=")" separator="OR">
                <if test="name == 'Jonathan' or name == 'Fuji'">
                    name = #{name}
                </if>
            </foreach>
        </where>
    </select>
    <insert id="testInsertMulti">
        INSERT INTO
        fruits
        (
        name,
        category,
        price
        )
        VALUES
        <foreach collection="fruits" item="fruit" separator=",">
            (
            #{fruit.name},
            #{fruit.category},
            ${fruit.price}
            )
        </foreach>
    </insert>
    <select id="testBind">
        <bind name="likeName" value="'%' + name + '%'"/>
        SELECT
        name,
        category,
        price
        FROM
        fruits
        WHERE
        name like #{likeName}
    </select>
</mapper>`,
		"SELECT `name`,`category`,`price` FROM `fruits` WHERE `category`=? AND `price`>?;\n"+
			"SELECT `name`,`category`,`price` FROM `fruits` WHERE `category`=?;\n"+
			"SELECT `name`,`category`,`price` FROM `fruits` WHERE 1=1 AND `category`=? AND `price`=? AND `name`=\"Fuji\";\n"+
			"SELECT `name`,`category`,`price` FROM `fruits` WHERE `category`=\"apple\" OR `price`=200;\n"+
			"SELECT `name`,`category`,`price` FROM `fruits` WHERE `category`=\"apple\" AND `price`=?;\n"+
			"UPDATE `fruits` SET `category`=?, `price`=? WHERE `name`=?;\n"+
			"SELECT `name`,`category`,`price` FROM `fruits` WHERE `name`=? AND `category`=? AND `price`=? AND `category`=\"apple\";\n"+
			"SELECT `name`,`category`,`price` FROM `fruits` WHERE `category`=\"apple\" AND (`name`=? OR `name`=?);\n"+
			"INSERT INTO `fruits` (`name`,`category`,`price`) VALUES (?,?,?),(?,?,?);\n"+
			"SELECT `name`,`category`,`price` FROM `fruits` WHERE `name` LIKE ?;",
	)
}
